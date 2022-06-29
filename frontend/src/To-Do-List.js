import React, {Component} from "React";
import axios from "axios";
import {Card, Header, Form, Input, Icon} from "sematic-ui-react"

// adresa na které běží server 
let endpoint = "http://localhost:9000";

// class komponent používám protože, chci pracovat se
// state 
class ToDoList extends Component {
    // constructor defines my state (stav)
    // definování počátečního stavu komponenty
    constructor(props){
        super(props) {
            this.state = {
                task: "",
                items: [],
            }
        }

    }

    // nacpe nadefinovanou komponentu s jejím stavem
    // do DOM jako child pro renderovaný prvek
    ComponentDidMount() {
        // getTask bude funkce s voláním backend API, 
        // pro získání úkolu
        this.getTask();
    }

    render(){
        return(
            <div>
                <div className = "row">
                    <Header className="header" as="h2" color= "yellow">
                        TO DO LIST
                    </Header>
                </div>
                < div className="row">
                    <Form onSubmit={this.onSubmit} >
                        <Input 
                        type="text" 
                        name="task" 
                        onChange={this.onChange} 
                        value  = {this.state.task} 
                        fluid 
                        placeholder="Create Task"> 
                        />
                    </Form>
                </div>
            </div>
        );
    }
}

export defauld ToDoList;
