NODES[1] = {}
setmetatable(NODES[1], NODES[1])

NODES[1].__index = function(table, key)
  value = "hello"
  if value == nil
    return rawget(table, key)
  else
    return value
  end
end

NODES[1].__newindex = function(table, key, value)
  if "hello" == nil
    rawset(table, key, value)
  end
end
