#Dir.entries("/storage1/osd/current").select {|f| File.directory? f}

#Dir.glob('*').select {|f| File.directory? f}

#Total number of osd
osd = `ceph osd stat`
puts "===========================> osd stat ============================="
puts osd
#Number of osd up

#number of replica
replica = `ceph osd dump | grep "replicated size"`

puts "===========================> osd replica ============================="
puts replica


Dir.glob("/storage1/osd/current/**/*").each do |f| 
  if File.directory?(f)
        file = f.split('/').last
        if file.include? "_head"
        hex_file = file.split('_').first
        puts `ceph pg map #{hex_file}`
        end
  end
end
