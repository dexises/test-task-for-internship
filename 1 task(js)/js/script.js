fetch('https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1')
  .then(response => response.json())
  .then(data => {
    const currencyTable = document.getElementById('currencyTable');
    
    data.forEach(currency => {
      const row = currencyTable.insertRow();
      
      if (row.rowIndex <= 5) {
        row.classList.add('highlight');
      }
      
      if (currency.symbol === 'usdt') {
        row.classList.add('usdt');
      }
      
      const idCell = row.insertCell();
      const symbolCell = row.insertCell();
      const nameCell = row.insertCell();
      
      idCell.textContent = currency.id;
      symbolCell.textContent = currency.symbol;
      nameCell.textContent = currency.name;
    });
  })
  .catch(error => console.error('Error:', error));
