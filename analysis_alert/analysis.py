import sqlite3
import pandas as pd

# ----------------------------
# Load raw data
# ----------------------------

with sqlite3.connect("../scraper/ram_prices.db") as conn:
    df = pd.read_sql_query("SELECT * FROM ram_prices", conn)

# normalize column safety (future-proofing)
if "MARKET" not in df.columns:
    df["MARKET"] = "unknown"

df["RAM_TYPE"] = df["RAM_TYPE"].fillna("unknown")
df["MARKET"] = df["MARKET"].fillna("unknown")

# remove garbage
df = df[df["RAM_TYPE"].str.lower() != "unknown"]


with sqlite3.connect("market_snapshot.db") as snapshot:

    snapshot.execute("""
    CREATE TABLE IF NOT EXISTS market_snapshots(
        id INTEGER PRIMARY KEY AUTOINCREMENT,

        market TEXT NOT NULL,
        ram_type TEXT NOT NULL,

        avg_price REAL NOT NULL,
        min_price INTEGER NOT NULL,
        max_price INTEGER NOT NULL,

        listing_count INTEGER NOT NULL,

        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
    """)

    groups = df.groupby(["MARKET", "RAM_TYPE"])

    for (market, ram_type), group in groups:

        if group.empty:
            continue

        avg_price = round(group["PRICE"].mean(), 2)
        min_price = int(group["PRICE"].min())
        max_price = int(group["PRICE"].max())
        listing_count = len(group)

        snapshot.execute(
            """
            INSERT INTO market_snapshots(
                market,
                ram_type,
                avg_price,
                min_price,
                max_price,
                listing_count
            )
            VALUES (?, ?, ?, ?, ?, ?)
            """,
            (
                market,
                ram_type,
                avg_price,
                min_price,
                max_price,
                listing_count,
            ),
        )

        print(
            f"[SNAPSHOT] {market:<8} {ram_type:<5} "
            f"avg={avg_price:,.0f} "
            f"min={min_price} "
            f"max={max_price} "
            f"count={listing_count}"
        )

    snapshot.commit()

print("[✓] Market snapshots created successfully")
