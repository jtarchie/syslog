-- define a destination to write messages to
-- this is a simple file type
-- each line will be a syslog message
destination('app123_logs', {
  type='file',
  config={
    path='/var/log/app123.log',
    message=function(message)
      -- add a 'log file' string to each app id
      message.app_id = 'log file ' + message.app_id
      return ({message=message})
    end
  }
})

-- listen on port 8080 for syslog messages
listen("udp", 65000, function(message) -- callback for the how to redirect the message and modify it
  if message.app_id == "123" then
    return({message=message, destination={'app123_logs'}})
  end

  return(nil)
end)