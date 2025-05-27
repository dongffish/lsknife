一、调用说明:
sqlite 可以指定数据库文件名，如 :
    db = sqlite.GormSQLite("mydb.data") 
或db = sqlite.GormSQLite("/home/data/db.dat") 
或db = sqlite.GormSQLite("d:\\db\abc.data")  

也可以不指定数据文件名，此时数据文件名默认放在 data\db.dat文件中,如：
sqlite.GormSQLite("")   

二、data目录
db_template.dat是空的sqlite3格式的文件。
db.dat也是空的sqlite3格式的文件。