import React from 'react';
import ReactDOM from 'react-dom';

import ReceiptForm from './component/receipt_form'

window.React = React;

window.onload = function () {
  ReactDOM.render(<ReceiptForm />, document.getElementById('app-container'));
};
