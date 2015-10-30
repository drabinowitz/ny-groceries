import React from 'react';
import request from '../ajax';
import DateInput from 'cmpnt/date/input';
import SelectDropdown from 'cmpnt/select/dropdown';

import ProductForm from './product_form';
import ProductTable from './product_table';
import PurchaseForm from './purchase_form';

export default class ReceiptForm extends React.Component {
  constructor() {
    super()
    this.state = {
      products: null,
      stores: null,
      selectedStore: null,
      productPurchases: [],
    };
  }

  onSubmitReceipt(e) {
    e.preventDefault();
    e.stopPropagation();
  }

  componentWillMount() {
    request('products/').then(data => this.setState({products: data}));
    request('stores/').then(data => this.setState({stores: data}));
  }

  submitNewProduct(product) {
    request('products/', {
      type: 'PUT',
      data: JSON.stringify(product),
    }).then(data => this.setState({
      products: this.state.products.concat([data]),
    }));
  }

  render() {
    if (!this.state.products || !this.state.stores) {
      return <h1>Loading App</h1>
    }
    let contentStores = this.state.stores.map(store => ({
      content: store.name,
      id: store.id,
    }));
    let productIds = this.state.productPurchases.map(p => p.product_id);
    let purchaseForms = this.state.productPurchases.map(p => (
      <PurchaseForm
        key={p.product_id}
        products={this.state.products}
        purchase={p}
        onPurchaseChange={this.onPurchaseChange.bind(this)}
        />
    ));
    return (
      <div>
        <h1>ReceiptForm</h1>
        <hr />
        <h4>Receipt Store</h4>
        <SelectDropdown
          selected={this.state.selectedStore}
          options={contentStores}
          onChange={this.onStoreDropdownChange.bind(this)}
          />
        <hr />
        <h4>Receipt Total</h4>
        <input
          type='number'
          value={this.state.receiptTotal}
          onChange={this.onReceiptTotalChange.bind(this)}
          />
        <hr />
        {purchaseForms}
        <hr />
        <button onClick={this.onSubmitReceipt}>Add Receipt</button>
        <hr />
        <ProductForm onSubmit={this.submitNewProduct.bind(this)} />
        <hr />
        <ProductTable
          onProductAdd={this.onProductPurchaseAdd.bind(this)}
          productIds={productIds}
          rows={this.state.products}
          />
      </div>
    );
  }

  onStoreDropdownChange(selectedStore) {
    this.setState({selectedStore});
  }

  onProductPurchaseAdd(product_id) {
    this.setState({
      productPurchases: this.state.productPurchases.concat([{product_id}]),
    });
  }

  onPurchaseChange(purchase) {
    let index = 0;
    let productPurchases = this.state.productPurchases.filter((p, i) => {
      if (p.product_id === purchase.product_id) {
        index = i;
        return false;
      }
      return true;
    });
    productPurchases.splice(index, 0, purchase);
    this.setState({productPurchases});
  }

  onReceiptTotalChange(e) {
    e.preventDefault();
    e.stopPropagation();
    this.setState({receiptTotal: e.target.value});
  }
};
