import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class CustomersRoute extends Route {
  @service store; // store injected

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  async model(params) {
    const customers = await this.store.query('customer', {
      page: params.page,
      per_page: 10,
    });
    return customers;
  }
}
