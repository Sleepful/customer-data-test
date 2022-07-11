import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class CustomersController extends Controller {
  queryParams = ['page'];
  @tracked page = 1;
  get prev() {
    const page = this.page;
    if (page == 1) return 1;
    return page - 1;
  }
  get next() {
    if (this.page == this.totalPages) {
      return this.page;
    }
    return this.page + 1;
  }
  get totalPages() {
    const page = this.page;
    const model = this.model;
    const perPage = model.meta.per_page;
    const totalRecords = model.meta.total;
    const totalPages = totalRecords / perPage;
    return totalPages;
  }
}
