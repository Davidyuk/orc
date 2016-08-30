import {RequestOptions, RequestOptionsArgs} from "@angular/http";
import {Injectable, Inject} from "@angular/core";

@Injectable()
export class AppRequestOptions extends RequestOptions {
  constructor(@Inject('apiBaseUrl') private apiBaseUrl:string) {
    super();
  }

  merge(options?:RequestOptionsArgs):RequestOptions {
    if (options.url.indexOf('api/') == 0)
      options.url = this.apiBaseUrl + options.url;

    return super.merge(options);
  }
}
