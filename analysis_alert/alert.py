import sqlite3
import pandas as pd
import requests
import os
from dotenv import load_dotenv

load_dotenv()

BOT_TOKEN = os.getenv("TELEGRAM_BOT_TOKEN")
CHAT_ID = os.getenv("TELEGRAM_CHAT_ID")

if not BOT_TOKEN or not CHAT_ID:
    raise ValueError("Missing TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID")

snapshot = sqlite3.connect("market_snapshot.db")


def send_telegram_message(message: str):
    url = f"https://api.telegram.org/bot{BOT_TOKEN}/sendMessage"

    response = requests.post(
        url,
        data={
            "chat_id": CHAT_ID,
            "text": message,
            "parse_mode": "Markdown",
        },
        timeout=10,
    )

    if response.status_code != 200:
        print("Failed to send message:")
        print(response.text)

def get_last_two_snapshots(market: str, ram_type: str):
    query = """
    SELECT *
    FROM market_snapshots
    WHERE market = ?
      AND ram_type = ?
    ORDER BY created_at DESC
    LIMIT 2
    """

    return pd.read_sql_query(
        query,
        snapshot,
        params=(market, ram_type),
    )


pairs = pd.read_sql_query(
    """
    SELECT DISTINCT market, ram_type
    FROM market_snapshots
    WHERE LOWER(ram_type) != 'unknown'
    """,
    snapshot,
)


report = []
report.append("📊 *RAM Market Report*\n")

for _, row in pairs.iterrows():

    market = row["market"]
    ram_type = row["ram_type"]

    df_snap = get_last_two_snapshots(market, ram_type)

    if len(df_snap) < 2:
        continue

    latest = df_snap.iloc[0]
    previous = df_snap.iloc[1]

    if previous["avg_price"] == 0:
        continue

    change = (
        (latest["avg_price"] - previous["avg_price"])
        / previous["avg_price"]
    ) * 100

    if change <= -5:
        emoji = "📉"
    elif change >= 5:
        emoji = "📈"
    else:
        emoji = "➖"

    report.append(
        (
            f"{emoji} *{market.upper()} | {ram_type.upper()}*\n"
            f"Avg: {latest['avg_price']:.0f} KES\n"
            f"Min: {latest['min_price']:.0f} KES\n"
            f"Max: {latest['max_price']:.0f} KES\n"
            f"Listings: {latest['listing_count']}\n"
            f"Change: {change:.2f}%\n"
        )
    )

message = "\n".join(report)
send_telegram_message(message)

snapshot.close()
