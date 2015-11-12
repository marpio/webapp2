import React from "react"
import ReactDOM from "react-dom"
import $ from "jquery"
import * as _ from "lodash"
import Item from "./item"

module.exports = class ShoppingCart extends React.Component {
  constructor(){
      super();
  }
    loadShoppingCartItems = () => {
      $.ajax({
        url: this.props.url,
        dataType: 'json',
        cache: false,
        success: function(data) {
          this.setState({
            order: data
          });
        }.bind(this),
        error: function(xhr, status, err) {
          console.error(this.props.url, status, err.toString());
        }.bind(this)
      });
    }
    onProductItemAmountChange = (amountChangeData) => {
      var o = Object.assign({}, this.state.order);
      var prodOrder = o.productOrders.filter(function(po) {
        return po.productID === amountChangeData.productID;
      })[0];
      prodOrder.amount = amountChangeData.amount;
      this.setState({
        order: o
      })
    }
    onProductItemDelete = (removeData) => {
      var o = Object.assign({}, this.state.order);
      _.remove(o.productOrders, function(po) {
        return po.productID === removeData.productID;
      });
      $.ajax({
        url: this.props.url,
        dataType: 'json',
        cache: false,
        data: JSON.stringify({
          productID: removeData.productID
        }),
        type: 'DELETE',
        contentType: 'application/json; charset=utf-8',
        success: function(result) {
          this.setState({
            order: result
          })
        }.bind(this)
      });
    }
    onConfirmOrder = () => {

    }
    componentDidMount = () => {
      this.loadShoppingCartItems();
    }
    state = {
      order: {
        productOrders: []
      },
    }
    render = () => {
      var self = this;
      var itemNodes = this.state.order.productOrders.map(function(item) {
        return ( <Item key = {
            item.productID
          }
          item = {
            item
          }
          onAmountChange = {
            self.onProductItemAmountChange
          }
          onItemDelete = {
            self.onProductItemDelete
          }
          />
        );
      });
      var totalPrice = 0;
      this.state.order.productOrders.forEach(function(po) {
        totalPrice += po.amount * po.productPrice;
      });
      var summary = this.state.order.productOrders.length > 0 ? < Summary totalPrice = {
        totalPrice
      }
      onConfirm = {
        this.onConfirmOrder
      }
      /> : null;
      return ( < ul >
        < li className = "row list-inline columnCaptions" >
        <span> &#13217;</span>
          <span>ITEM</span>
          <span>Price</span>
        </li>
        {itemNodes}
        {summary}
      </ul>
    );
  }
}


      // var ServiceProvider = React.createClass({
      //   render: function() {
      //     return (
      //
      //     );
      //   }
      // });

      var Summary = React.createClass({
          handleConfirm: function() {
            this.props.onConfirm();
          },
          render: function() {
            return ( < li className = "row totals" >
              < span className = "itemName" > Total: < /span> < span className = "price" > EUR {
              this.props.totalPrice
            } < /span> < span className = "order" > < button type = "submit"
            className = "btn btn-success"
            onClick = {
              this.handleConfirm
            } > Confirm < /button> < /span > < /li>
          );
        }
      });
