import { loadMoonPay } from '@moonpay/moonpay-js';

// Browser
const moonPay = await loadMoonPay();
const moonPaySdk = moonPay({
    flow: 'sell',
    environment: 'sandbox',
    variant: 'overlay',
    params: {
        apiKey: 'pk_test_paQPDUbJpqLFsibzdhBEX0Y7SsAmFkIL',
        theme: 'dark',
        quoteCurrencyCode: 'usd',
        baseCurrencyAmount: '.01',
        defaultBaseCurrencyCode: 'eth'
    }
});

moonPaySdk.show()