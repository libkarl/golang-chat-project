import React from "react";
import "./App.css";
import {Container} from "sematic-ui-react";
import ToDoList from "./ToDoList";

function App() {
  return(
    <div>
      <Container>
        <ToDoList/>
      </Container>
    </div>
  );
  
}

export default App;