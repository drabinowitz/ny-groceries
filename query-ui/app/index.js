import React from 'react';
import ReactDOM from 'react-dom';

import ProductCostsTable from './component/product_costs_table'

window.React = React;

window.onload = function () {
  ReactDOM.render(<ProductCostsTable />, document.getElementById('app-container'));
};

