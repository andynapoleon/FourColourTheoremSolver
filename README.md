# 🗺️ Four-Color Map Theorem Solver 🎨

## 🌟 Summary

Welcome to the Four-Color Map Theorem Solver, an interactive web application demonstrating one of mathematics' most significant theorems in graph theory and topology. This tool provides a practical implementation of the Four-Color Theorem, which states that any planar map can be colored using no more than four colors while ensuring no adjacent regions share the same color, through an intuitive, user-friendly interface.

Initially developed collaboratively during the **DevelopEd 2.0 Hackathon (2023)** 🏆 and further advanced 🚀. Our implementation allows users to:
- Create custom maps through an interactive canvas
- Automatically generate mathematically valid four-color solutions
- Save and download colored maps
- Visualize the theorem's practical applications

## 🚀 Tech Stack & Architectural Design

1. Front-End: **Next.js** 13 App Router (with **React.js** and **TailwindCSS**), deployed on **Vercel**.
2. Back-End: Microservices Architecture with **Go**, **Python Flask**, **gRPC**, **RabbitMQ**, **MongoDB**, and **PostgreSQL**.
3. DevOps: Containerized services using **Docker** and orchestrated deployments with **Kubernetes** (on **AWS EKS**), ensuring high availability and scalability.
4. Services: API gateway, authentication, map solver (map coloring), map storage, and logger, communicating with each other via different protocols such as **REST**, **gRPC**, and **AMQP**.
5. Constraint-satisfaction algorithm implemented with Answer-Set Programming (AST) in **Python** & **clingo** used to solve the four-color map theorem.

![map-solver drawio(5)](https://github.com/user-attachments/assets/042ca7d4-652a-4d50-b37c-bc19dd18171a)

## 📞 Contact

- Andy Tran ([anhquoctran006@gmail.com](mailto:anhquoctran006@gmail.com))
- Riley Kinahan ([rdkinaha@ualberta.ca](mailto:rdkinaha@ualberta.ca))
