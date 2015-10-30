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
  render() {
    let cost = this.props.purchase.cost || 0;
    let quantity = this.props.purchase.quantity || 1;
    let unit = this.props.purchase.unit || 'unit';
    let product = this.props.products.filter(p =>(
      p.id === this.props.purchase.product_id
    ))[0]
    return (
      <div>
        <h4>{product.category} {product.sub_category}</h4>
        <span> cost: </span>
        <input
          type='number'
          value={cost}
          onChange={this.onCostChange.bind(this)}
          />
        <span> quantity: </span>
        <input
          type='number'
          value={quantity}
          onChange={this.onQuantityChange.bind(this)}
          />
        <span> unit: </span>
        <SelectDropdown
          options={options}
          selected={options.filter(o => o.id === unit)[0]}
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
    let quantity = e.target.value;
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
    var purchase = {};
    for (let key in this.props.purchase) {
      purchase[key] = this.props.purchase[key];
    }
    purchase.cost = cost;
    this.props.onPurchaseChange(purchase);
  }
}

PurchaseForm.propTypes = {
  products: React.PropTypes.array.isRequired,
  purchase: React.PropTypes.object.isRequired,
  onPurchaseChange: React.PropTypes.func.isRequired,
};

export default PurchaseForm;
