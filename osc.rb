#!/usr/bin/env ruby

require 'osc-ruby'

Server = Struct.new(:ip, :port, :path, :dim, :valid)


servers = {
  strip: Server.new('10.100.7.28', 13000, '/led-strip/', [1, 192, 6], /^[0123456789abcdef]+$/),
  board: Server.new('10.100.7.28',12000, '/', [38, 96, 1], /^[01]+$/)
}

server = servers[ARGV[0] ? ARGV[0].to_sym : nil]
if server.nil?
  STDERR.puts "invalid server: '#{ARGV.join(',')}'"
  STDERR.puts "please enter one of #{servers.keys.map(&:to_s).join('|')}"
  exit 1
end

client = OSC::Client.new(server.ip, server.port)
line_len = server.dim.reduce(:*)

while line = STDIN.gets
  line = line.strip
  err = false
  if line.size != line_len
    STDERR.puts "Line is not the appropriate length (#{line.size} vs #{line_len})"
    err = true
  end

  unless server.valid.match(line)
    STDERR.puts "Line contains invalid characters"
    err = true
  end

  if err
    STDERR.puts "\"#{line[0..60]}...\""
  else
    msg = OSC::Message.new(server.path, line)
    client.send(msg)
  end
end



