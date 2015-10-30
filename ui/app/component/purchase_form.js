import React from 'react';
import SelectDropdown from 'cmpnt/select/dropdown';

let options = [
  {id: 'unit', content: 'unit'},
  {id: 'qt', content: 'qt'},
  {id: 'pt', content: 'pt'},
  {id: 'oz', content: 'oz'},
  {id: 'lb', content: 'oz'},
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
      unit: props.purchase.unit || 'unit',
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
          />
        <span> unit: </span>
        <SelectDropdown
          options={options}
          selected={options.filter(o => o.id === this.state.unit)[0]}
          onChange={this.onUnitChange.bind(this)}
          />
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
    let quantity = parseInt(e.target.value);
    e.preventDefault();
    e.stopPropagation();
    var purchase = {};
    for (let key in this.props.purchase) {
      purchase[key] = this.props.purchase[key];
    }
    purchase.quantity = quantity;
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
}

PurchaseForm.propTypes = {
  products: React.PropTypes.array.isRequired,
  purchase: React.PropTypes.object.isRequired,
  onPurchaseChange: React.PropTypes.func.isRequired,
};

export default PurchaseForm;
