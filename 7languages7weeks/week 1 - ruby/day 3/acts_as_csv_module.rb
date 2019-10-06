module ActsAsCsv
	def self.included(base)
		base.extend ClassMethods
	end
	
	module ClassMethods
		def acts_as_csv
			include InstanceMethods
		end
	end
	
	module InstanceMethods
		attr_accessor :headers, :csv_contents
		
		def initialize
			read
		end
		
		def read
			@csv_contents = []
			file = File.new(self.class.to_s.downcase + '.csv')
			@headers = file.gets.chomp.split(', ')
			
			file.each do |row|
				@csv_contents << CsvRow2.new(row.chomp.split(', '), @headers)
			end
		end

		
		def each
			@csv_contents.each { |row| yield row }
		end
	end
	
	class CsvRow
		attr_accessor :headers, :row_contents
		
		def initialize(row, headers)
			@row_contents = row
			@headers = headers
		end
		
		def method_missing(name, *args, &block)
			index = @headers.index(name.to_s)
			if index
				@row[index]
			else
				super
			end 
		end
	end
	
	class CsvRow2
		attr_accessor :row_contents
		
		def initialize(row, headers)
			@row_contents = Hash[headers.zip(row)]
		end
		
		def method_missing(name, *args, &block)
			column = @row_contents[name.to_s]
			if column
				column
			else
				super
			end 
		end
	end
end

class RubyCsv
	include ActsAsCsv
	acts_as_csv
end

m=RubyCsv.new
puts m.headers.inspect
puts m.csv_contents.inspect
puts '----'
m.each{|row| puts row.one}