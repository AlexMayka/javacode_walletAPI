counter = 0

request = function()
  counter = counter + 1
  if counter % 2 == 0 then
    local operation = "DEPOSIT"
    local body = string.format(
      '{"walletId":"1c63a43f-aacd-47b0-bc3b-535e69c6ed4c","operationType":"%s","amount":100}',
      operation
    )
    return wrk.format("POST", "/api/v1/wallet", {["Content-Type"] = "application/json"}, body)
  else
    return wrk.format("GET", "/api/v1/wallets/1c63a43f-aacd-47b0-bc3b-535e69c6ed4c")
  end
end