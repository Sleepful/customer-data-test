import Model from '@ember-data/model';
import { attr } from '@ember-data/model';

export default class CustomerModel extends Model {
  //@attr id;
  @attr attributes;
  @attr events;
  @attr last_updated;
  get fullName() {
    return `${this.attributes?.first_name} ${this.attributes?.last_name}`;
  }
}
