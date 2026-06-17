read MMC cid objects and make em human readable

#### MMC CID table

| Field | Description |
| :--: | :-- |
| MID | 8 - bit Manufacturer ID (assigned by SD-3C LLC) |
| OID | 16 - bit ASCII OEM/Application ID |
| PNM | 40 - bit ASCII Product Name |
| PRV | 8 - bit Product Revision in BCD (e.g., 2.0 = 0x20) |
| Reserved | 4 - bit |
| PSN | 32 - bit Product Serial Number |
| MDT | 12 - bit Manufacturing Date in YYM format (offset from year 2000) |
| CRC | 7 - bit CRC checksum for validation |
| N/A | 1- bit Not used, always 1 | 
