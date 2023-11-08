package model

const q_NEWBOAT string = "INSERT INTO boat(id, name, type_id) values(?, ?, ?)"
const q_GETBOATS string = "SELECT * from boat"
const q_GETBOAT string = "SELECT * from boat where id=?"
const q_FINDBOAT string = "SELECT * from boat where name=?"
const q_LASTBOATSTATE string = "SELECT * from boat_state where boat_id=? order by timestamp desc limit 1"
const q_NEWBOATTYPE string = "INSERT INTO boat_type(id, name) values(?, ?)"
const q_NEWSAILTYPE string = "INSERT INTO sail(id, name) values(?, ?)"
const q_NEWSTATE string = `INSERT INTO boat_state
		(timestamp, boat_id, voyage_id, latitude, longitude, sog, cog, spd, hdg, awa, aws, twa, tws, 
		divedegrees, drift, foilleft, foilright, heeldegrees, keelangle, waterballast, weatherhelm, sails) 
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
