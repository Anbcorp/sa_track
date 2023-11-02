package model

const q_NEWBOAT string = "INSERT INTO boat(id, name, type_id) values(?, ?, ?)"
const q_GETBOATS string = "SELECT * from boat"
const q_GETBOAT string = "SELECT * from boat where id=?"
const q_FINDBOAT string = "SELECT * from boat where name=?"
