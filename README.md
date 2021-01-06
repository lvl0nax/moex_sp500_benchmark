# Compare your stocks portfolio with MOEX sp500 becnhmark (from Sber)

You have to provide XLSX file (broker_report.xlsx) which contains the following columns - 
HEADERS (must have exactly the same name):
 - DATE - transaction date
 - PRICE - spent money by transaction
 - CURRENCY - transaction currency (app works with RUB transactions only)
 - you can have another columns, the app will ignore them.
 
 
## FLow
1. parse provided XLSX document
2. fetch sp500 benchmark prices with dates of parsed transactions
3. create a new XLSX document (tmp_result.xlsx) with the number of sp500 benchmark items for each transaction. 

You can then add up all the benchmark items and multiply by the last benchmark price. This way, you would have got the possible price of the portfolio if you had been buying only benchmarks items.

