# Athens - Data sync pipeline

## 📑 Objective

The objectives of this task are as follows:

- To build a data pipeline that syncs and transforms the data stored in a PostgresSQL Database to Elasticsearch, making the data searchable in a fuzzy and relational manner.
- To build a REST API service that will query Elasticsearch.

## 🛠 Requirements

- Create a Postgres Database as defined in the [Database Schema Section](https://www.notion.so/Backend-Take-home-task-1de25b809b594d7fb8cb473a4d234d28?pvs=21) and add some seed data.
- The data pipeline should be able to handle CRUD operations made to the database and sync them to Elasticsearch.
- Create RESTful search APIs in the language of your choice (**preferably GoLang**) that will query Elasticsearch for the following functionality
    1. **Search for projects created by a particular user**
        - This API should return the project details, along with the hashtags used in the project, and the details of the users that have created the project.
    2. **Search for projects that use specific hashtags**
        - This API should return the project details, along with the hashtags used in the project, and the details of the users that have created the project.
    3. **Full-text fuzzy search for projects**
        - This API should allow fuzzy searching using the project `slug` and the `description`.
        - This API should return the `hashtags` used by the projects along with details of the users that have created the project.
- A Terraform/IaC script that provisions the required resources on the cloud of your choice (Preferably AWS).

## 💽 Database Schema

Here is a DB Diagram schema that can be used as a reference, feel free to make your own assumptions.

![Untitled.png](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/1ca65691-cb8b-4154-bc16-099fd7084925/Untitled.png)

- DB Diagram Snippet
    
    ```sql
    Table users as U {
      id int [pk, increment]
      name varchar
      created_at timestamp
    }
    
    Table hashtags as H {
      id int [pk, increment]
      name varchar
      created_at timestamp
    }
    
    Table projects as P {
      id int [pk, increment]
      name varchar
      slug varchar
      description text
      created_at timestamp
    }
    
    Table project_hashtags {
      hashtag_id int [ref: > hashtags.id]
      project_id int [ref: > projects.id]
    }
    
    Table user_projects {
      project_id int [ref: > projects.id]
      user_id int [ref: > users.id]
    }
    ```