import React from 'react';

class FoundPurchasesModal extends React.Component {

  render() {
    const purchaseRows = this.props.purchases.map(purchase => {
      return (
        <div key={purchase.id}>
          <p>quantity: {purchase.quantity} ({purchase.unit})</p>
          <p>cost: {purchase.cost}</p>
          <button onClick={this.props.onPurchaseClick.bind(null, purchase)}>
            Add
          </button>
          <hr />
        </div>
      );
    });
    return (
      <div>
        <h2>Found Purchases Modal</h2>
        <button onClick={this.props.onCancelClick}>Cancel</button>
        {purchaseRows}
        <hr />
        <hr />
        <hr />
        <hr />
        <hr />
        <hr />
        <hr />
        <hr />
      </div>
    );
  }
}

FoundPurchasesModal.propTypes = {
    purchases: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
    onCancelClick: React.PropTypes.func.isRequired,
    onPurchaseClick: React.PropTypes.func.isRequired
};

export default FoundPurchasesModal
