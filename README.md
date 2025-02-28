# Schedule_Table
project automate schedule resource management

implement by golang programming use gin framework

framework : gin, grom
front-end : next, ant design
database  : postgresql16

![screenshot](doc/Screenshot%20from%202025-02-04%2019-49-46.png)
![screenshot](doc/Screenshot%20from%202025-02-04%2019-49-37.png)
![screenshot](doc/Screenshot%20from%202025-02-04%2019-49-30.png)

*** design demo application

ถ้าปัญหามีเงื่อนไขเชิงเส้น → Linear Programming (LP), Integer Programming (IP), MIP
ถ้าปัญหามีฟังก์ชันไม่เป็นเชิงเส้น → Non-Linear Programming (NLP), Convex Optimization
ถ้าต้องการใช้วิธีแบบอัลกอริทึมเฮอริสติก → Genetic Algorithm (GA), Simulated Annealing (SA), PSO
ถ้ามีโครงสร้างซ้ำซ้อน → Dynamic Programming (DP), Constraint Programming (CP)

นอกจาก **Linear Programming (LP)** แล้ว ยังมีเทคนิคการทำ **Optimization** อื่น ๆ ที่ใช้ได้ในหลายบริบท ขึ้นอยู่กับลักษณะของปัญหาและข้อจำกัดที่เกี่ยวข้อง นี่คือตัวอย่างของเทคนิคสำคัญที่ใช้กัน:

---

## **1. Non-Linear Programming (NLP)**
เมื่อฟังก์ชันเป้าหมาย (Objective Function) หรือข้อจำกัด (Constraints) มีความไม่เป็นเชิงเส้น  
🔹 วิธีที่ใช้:  
- **Gradient Descent** (เหมาะกับ Machine Learning)  
- **Newton’s Method** (ใช้อนุพันธ์ระดับที่สอง)  
- **Conjugate Gradient Method**  

---

## **2. Integer Programming (IP)**
เหมือน **Linear Programming** แต่ตัวแปรต้องเป็นจำนวนเต็ม  
🔹 วิธีที่ใช้:  
- **Branch and Bound**  
- **Cutting Plane Method**  
- **Branch and Cut**  

**📌 ตัวอย่าง:**  
ปัญหา **Knapsack Problem**, **Scheduling Problem**, **Vehicle Routing Problem**  

---

## **3. Mixed-Integer Programming (MIP)**
เป็นการรวมกันระหว่าง **Linear Programming (LP)** และ **Integer Programming (IP)**  
🔹 ใช้เมื่อบางตัวแปรเป็นจำนวนเต็ม และบางตัวแปรเป็นค่าต่อเนื่อง  
🔹 นิยมใช้ใน **โลจิสติกส์, การจัดสรรทรัพยากร, การวางแผนกำลังคน**  

---

## **4. Dynamic Programming (DP)**
ใช้แนวคิด **แบ่งปัญหาใหญ่เป็นปัญหาเล็ก ๆ และแก้ซ้ำ**  
🔹 เหมาะกับปัญหาที่มี **Overlapping Subproblems** และ **Optimal Substructure**  
🔹 นิยมใช้ใน **Algorithm Design** เช่น **Fibonacci, Shortest Path (Dijkstra), Longest Common Subsequence (LCS)**  

---

## **5. Genetic Algorithm (GA)**
🔹 เลียนแบบหลักการคัดเลือกทางธรรมชาติ  
🔹 ใช้ในปัญหาที่ไม่มีฟังก์ชันเชิงเส้นชัดเจน เช่น **Neural Network Optimization, Traveling Salesman Problem (TSP)**  
🔹 ขั้นตอน: **Selection → Crossover → Mutation → Evaluation**  

---

## **6. Simulated Annealing (SA)**
🔹 ใช้หลักการเดียวกับการทำให้โลหะเย็นลง (Annealing)  
🔹 ใช้ลดการติดอยู่ใน Local Optimum  
🔹 เหมาะกับปัญหา Combinatorial เช่น **TSP, Job Scheduling**  

---

## **7. Particle Swarm Optimization (PSO)**
🔹 จำลองพฤติกรรมของฝูงนกหรือปลาในการหาจุดที่ดีที่สุด  
🔹 ใช้หาค่าที่เหมาะสมในปัญหา Continuous Optimization  

---

## **8. Convex Optimization**
🔹 ใช้กับฟังก์ชันที่เป็น **Convex Function** ซึ่งรับประกันว่าจุดต่ำสุดที่หาได้เป็น **Global Optimum**  
🔹 ใช้ใน **Machine Learning, Finance, Signal Processing**  

---

## **9. Constraint Programming (CP)**
🔹 ใช้แก้ปัญหาที่มีเงื่อนไขซับซ้อน เช่น **Sudoku, Scheduling, Routing**  
🔹 ใช้ **Constraint Propagation** เพื่อตัดตัวเลือกที่เป็นไปไม่ได้ออก  

---

### 🔥 **สรุป**
ถ้าปัญหามีเงื่อนไขเชิงเส้น → **Linear Programming (LP), Integer Programming (IP), MIP**  
ถ้าปัญหามีฟังก์ชันไม่เป็นเชิงเส้น → **Non-Linear Programming (NLP), Convex Optimization**  
ถ้าต้องการใช้วิธีแบบอัลกอริทึมเฮอริสติก → **Genetic Algorithm (GA), Simulated Annealing (SA), PSO**  
ถ้ามีโครงสร้างซ้ำซ้อน → **Dynamic Programming (DP), Constraint Programming (CP)**  

---

อยากโฟกัสที่ปัญหาหรืออัลกอริทึมแบบไหนเป็นพิเศษไหม? 😊