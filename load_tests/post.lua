counter = 0

request = function()
  counter = counter + 1
  local operationType = (counter % 2 == 0) and "DEPOSIT" or "WITHDRAW"
  local amount = 100

  local body = '{"walletId":"1c63a43f-aacd-47b0-bc3b-535e69c6ed4c","operationType":"'..operationType..'","amount":'..amount..'}'
  return wrk.format("POST", "/api/v1/wallet", {["Content-Type"]="application/json"}, body)
end