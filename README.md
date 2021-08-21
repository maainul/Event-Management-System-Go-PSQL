# Event Management System

How to Start :

1. Create a database in the postgres

```
    	dbParams := " " + "user=postgres"
        dbParams += " " + "host=localhost"
        dbParams += " " + "port=5432"
        dbParams += " " + "dbname=dbevent"
        dbParams += " " + "password=password"
```

2.  Go to the project directory and run :

        go run migrations/migrate.go up

3.  Run cmd : to the project root directory

     ~/go/src/github.com/Event-Management-System-Go-PSQL

    go run main.go

    http://localhost:8080/

4.  Project Structure :

![projectstruct](https://user-images.githubusercontent.com/37740006/130315655-0d5c5442-0549-4802-88c6-6cf555b00aad.png)

5.  Project urls:

        indexpage : http://localhost:8080/

    ![Screenshot from 2021-08-21 14-25-06](https://user-images.githubusercontent.com/85335954/130315830-6acc4961-4a18-49b9-b8e6-0e68773738a1.png)

        user registration : http://localhost:8080/user/create

    ![userreh](https://user-images.githubusercontent.com/85335954/130315937-8197c169-cd94-412c-8fa6-8b91027da5ce.png)

        login :  user login : http://localhost:8080/login

    ![login](https://user-images.githubusercontent.com/85335954/130315976-214aa843-d656-4227-81e2-b6c0c9cb4c6d.png)

go to postgres ui and change is_admin : true and login again

    admin-dashboard : http://localhost:8080/auth/admin-home

![admin-dash](https://user-images.githubusercontent.com/85335954/130316145-df95185f-6fad-43a1-9c44-48711a6196d4.png)

After that user can create :

    1. Event type
    2. Event
    3. Create Speaker
    4. Update Speaker

After that user can see :

    1. Event type list
    2. Event lsit
    3. Create Speaker list
    4. Feedback list

## Requirement

## Admin Functionality :

    1. Admin login
    2. Amin will create Event Type.Event type
    3. Admin can see event list, update and delete event list(CRUD of EVENT list)
    4. It allows the user to select from a list of event types.
    5. Add Venue
    6. New Booking with status
    7. can change staus

## User Functionality:

    1. User Can login
    2. User can signup(UseId,Name,Address,MobileNo,email,Password,ConfirmPass)
    3. User can see their details
    4. User can see venu list according to venue type along with cost
    5. View Venu form
    6. User can see their event list
    4. User Can Add  user to select the date and time of event, place and the event equipment’s.
    5.  user is given a receipt number for his booking.

# Database Tables

users

events

event-type

feedback

speaker

booking
