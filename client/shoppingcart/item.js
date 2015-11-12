import React from "react"
import ReactDOM from "react-dom"
import $ from "jquery"
import * as _ from "lodash"

module.exports = class Item extends React.Component {
  constructor() {
    super();
  }
  handleChange = (event) => {
    if (event.target.value > 0) {
      this.props.onAmountChange({
        productID: this.props.item.productID,
        amount: event.target.value
      });
    }
  }
  handleDelete = (event) => {
    this.props.onItemDelete({
      productID: this.props.item.productID
    });
  }
            render=()=> {
              var orderItemPrice = this.props.item.productPrice * this.props.item.amount;
              return ( < li className = "row" >
                <input type = "number" className = "quantity" value = {this.props.item.amount}  onChange = {this.handleChange}/>
                < span className = "text itemName" > {this.props.item.productName} < /span>
              < button type = "button" className = "remove btn btn-danger" onClick = { this.handleDelete } >
              < span className = "glyphicon glyphicon-trash" > < /span></button >
                < span className = "text price" > { orderItemPrice } < /span> < /li >
            );
          }
        }
