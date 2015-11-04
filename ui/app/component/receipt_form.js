import React from 'react';
import request from '../ajax';
import SelectDropdown from 'cmpnt/select/dropdown';
import moment from 'moment'

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
      date: moment().format('MM/DD/YYYY'),
      receiptTotal: 0,
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

  onSubmitReceipt() {
    request('receipt_uploads/', {
      type: 'PUT',
      data: JSON.stringify({
        receipt: {
          store_id: this.state.selectedStore.id,
          total: parseFloat(this.state.receiptTotal),
          date: this.state.date,
        },
        purchases: this.state.productPurchases,
      }),
    }).then(console.log.bind(console));
  }

  submitNewProduct(product) {
    request('products/', {
      type: 'PUT',
      data: JSON.stringify(product),
    }).then(data => this.setState({
      products: this.state.products.concat([data]),
      productPurchases: this.state.productPurchases.concat([{
        product_id: data.id
      }]),
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
        removePurchase={this.removePurchase.bind(this)}
        />
    ));
    let wrapperDiv = {
      display: 'inline-block',
      height:  '100%',
      width:  '50%',
      boxSizing: 'border-box',
      border: '1px solid black',
      float: 'left'
    }

    return (
      <div>
        <div className='left' style={wrapperDiv}>
          <h1>ReceiptForm</h1>
          <hr />
          <h4>Receipt Store</h4>
          <SelectDropdown
            selected={this.state.selectedStore}
            options={contentStores}
            onChange={this.onStoreDropdownChange.bind(this)}
          />
          <hr />
          <h4>Receipt Date</h4>
          <input
            value={this.state.date}
            onChange={this.onDateChange.bind(this)}
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
          <button onClick={this.onSubmitReceipt.bind(this)}>Add Receipt</button>
        </div>
        <div className='right' style={wrapperDiv}>
          <ProductForm onSubmit={this.submitNewProduct.bind(this)} />
          <hr />
          <ProductTable
            onProductAdd={this.onProductPurchaseAdd.bind(this)}
            productIds={productIds}
            rows={this.state.products}
          />
        </div>
      </div>
    );
  }

  onDateChange(e) {
    let date = e.target.value;
    e.stopPropagation();
    e.preventDefault();
    this.setState({date});
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

  removePurchase(purchase) {
    this.setState({
      productPurchases: this.state.productPurchases.filter(p => {
        return !(p.product_id === purchase.product_id);
      }),
    })
  }

  onReceiptTotalChange(e) {
    e.preventDefault();
    e.stopPropagation();
    this.setState({receiptTotal: e.target.value});
  }
};
