import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class CustomerRoute extends Route {
  @service store; // store injected

  async model({ customer_id }) {
    const customer = await this.store.findRecord('customer', customer_id);
    console.log(Object.entries(customer.attributes));
    return Object.entries(customer.attributes);
  }
}
