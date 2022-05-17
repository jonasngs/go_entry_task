user_id = 1

request = function()
    username = "user-" .. tostring(user_id)
    path = "/login" 
    wrk.method = "POST"
    wrk.body   = string.format("username=%s&password=password2", username)
    wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
    user_id = user_id + 1
    return wrk.format(nil, path)
end