# The Web Connected Chess Clock
An arduino based Chess Clock that allows for web based configuration


# API Endpoints Explorer

- api:/timecontrol/preset [GET]
<br>Retrieves all time control presets
- api:/timecontrol/preset/create [POST]
<br>Creates a time control preset
- api:/timecontrol/preset/{ID} [GET]
<br>Retrieve a specific preset
- api:/timecontrol/preset/{ID}/update [PATCH]
<br>Updates a preset
- api:/timecontrol/preset/{ID}/delete [DELETE]
<br>Used to delete a preset
- api:/arduino [GET]
<br>Gets the current device settings
- api:/arduino/display/preset [GET]
<br>Gets the current display preset
- api:/arduino/display/preset/create [POST]
<br>Creates a new display preset on the database
- api:/arduino/display/preset/{ID}/update [PATCH]
<br>Updates an existing preset
- api:/arduino/display/preset/{ID}/delete [DELETE]
<br>Deletes an existing preset

