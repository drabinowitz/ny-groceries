import React from 'react';
import SelectDropdown from 'cmpnt/select/dropdown';

let options = [
  {id: 'unit', content: 'unit'},
  {id: 'qt', content: 'qt'},
  {id: 'pt', content: 'pt'},
  {id: 'oz', content: 'oz'},
  {id: 'lb', content: 'lb'},
];

class PurchaseForm extends React.Component {
  constructor(props) {
    super(props)
    this.state = this._stateFromProps(this.props);
  }

  _stateFromProps(props) {
    return {
      cost: props.purchase.cost || 0,
      quantity: props.purchase.quantity || 1,
      unit: props.purchase.unit || 'oz',
      product: props.products.filter(p =>(
        p.id === props.purchase.product_id
      ))[0],
    };
  }

  componentWillReceiveProps(props) {
    this.setState(this._stateFromProps(props));
  }

  render() {
    return (
      <div>
        <h4>{this.state.product.category} {this.state.product.sub_category}</h4>
        <span> cost: </span>
        <input
          type='number'
          value={this.state.cost}
          onChange={this.onCostChange.bind(this)}
          onBlur={this.commitCostChange.bind(this)}
          />
        <span> quantity: </span>
        <input
          type='number'
          value={this.state.quantity}
          onChange={this.onQuantityChange.bind(this)}
          onBlur={this.commitQuantityChange.bind(this)}
          />
        <span> unit: </span>
        <SelectDropdown
          options={options}
          selected={options.filter(o => o.id === this.state.unit)[0]}
          onChange={this.onUnitChange.bind(this)}
          />
        <button onClick={this.removePurchase.bind(this)}>Delete</button>
      </div>
    );
  }

  onUnitChange(unit) {
    var purchase = {};
    for (let key in this.props.purchase) {
      purchase[key] = this.props.purchase[key];
    }
    purchase.unit = unit.id;
    this.props.onPurchaseChange(purchase);
  }

  onQuantityChange(e) {
    let quantity = e.target.value;
    e.preventDefault();
    e.stopPropagation();
    this.setState({quantity});
  }

  commitQuantityChange(e) {
    var purchase = {};
    for (let key in this.props.purchase) {
      purchase[key] = this.props.purchase[key];
    }
    purchase.quantity = parseFloat(this.state.quantity);
    this.props.onPurchaseChange(purchase);
  }

  onCostChange(e) {
    let cost = e.target.value;
    e.preventDefault();
    e.stopPropagation();
    this.setState({cost});
  }

  commitCostChange(e) {
    var purchase = {};
    for (let key in this.props.purchase) {
      purchase[key] = this.props.purchase[key];
    }
    purchase.cost = parseFloat(this.state.cost);
    this.props.onPurchaseChange(purchase);
  }

  removePurchase(e) {
    e.preventDefault();
    e.stopPropagation();
    this.props.removePurchase(this.props.purchase);
  }
}

PurchaseForm.propTypes = {
  products: React.PropTypes.array.isRequired,
  purchase: React.PropTypes.object.isRequired,
  onPurchaseChange: React.PropTypes.func.isRequired,
  removePurchase: React.PropTypes.func.isRequired,
};

export default PurchaseForm;
