import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class CustomerRoute extends Route {
  @service store; // store injected

  async model({ customer_id }) {
    const customer = await this.store.findRecord('customer', customer_id);
    console.log(customer.attributes);
    console.log(customer.created_at);
    customer.attributes.created_at = new Date(
      parseInt(customer.attributes.created_at)
    ).toISOString();
    return [
      ...Object.entries(customer.attributes),
      ['Last updated', new Date(customer.last_updated).toISOString()],
    ];
  }
}
