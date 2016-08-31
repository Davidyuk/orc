import { Injectable }    from '@angular/core';
import { Headers, Response, Http } from '@angular/http';

import 'rxjs/add/operator/toPromise';

import { Event } from './event';

@Injectable()
export class EventService {

  private eventsUrl = 'api/event';

  constructor(private http: Http) { }

  search(term: string) {
    return this.http
      .get(`${this.eventsUrl}?name=${term}`)
      .map((r: Response) => r.json() as Event[]);
  }

  getEvents() {
    return this.http.get(this.eventsUrl)
               .toPromise()
               .then(response => response.json() as Event[])
               .catch(this.handleError);
  }

  getEvent(id: number) {
    return this.getEvents()
               .then(events => events.find(event => event.id === id));
  }

  save(event: Event): Promise<Event>  {
    if (event.id) {
      return this.put(event);
    }
    return this.post(event);
  }

  delete(event: Event) {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');

    let url = `${this.eventsUrl}/${event.id}`;

    return this.http
               .delete(url, {headers: headers})
               .toPromise()
               .catch(this.handleError);
  }

  // Add new Event
  private post(event: Event): Promise<Event> {
    let headers = new Headers({
      'Content-Type': 'application/json'});

    return this.http
               .post(this.eventsUrl, JSON.stringify(event), {headers: headers})
               .toPromise()
               .then(res => res.json())
               .catch(this.handleError);
  }

  // Update existing Event
  private put(event: Event) {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');

    let url = `${this.eventsUrl}/${event.id}`;

    return this.http
               .put(url, JSON.stringify(event), {headers: headers})
               .toPromise()
               .then(() => event)
               .catch(this.handleError);
  }

  private handleError(error: any) {
    console.error('An error occurred', error);
    return Promise.reject(error.message || error);
  }
}
