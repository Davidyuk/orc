import { Injectable }    from '@angular/core';
import { Headers, Response, Http } from '@angular/http';

import 'rxjs/add/operator/toPromise';

import { Event } from './hero';

@Injectable()
export class EventService {

  private heroesUrl = 'api/event';

  constructor(private http: Http) { }

  search(term: string) {
    return this.http
      .get(`${this.heroesUrl}?name=${term}`)
      .map((r: Response) => r.json() as Event[]);
  }

  getHeroes() {
    return this.http.get(this.heroesUrl)
               .toPromise()
               .then(response => response.json() as Event[])
               .catch(this.handleError);
  }

  getHero(id: number) {
    return this.getHeroes()
               .then(heroes => heroes.find(hero => hero.id === id));
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

    let url = `${this.heroesUrl}/${event.id}`;

    return this.http
               .delete(url, {headers: headers})
               .toPromise()
               .catch(this.handleError);
  }

  // Add new Event
  private post(hero: Event): Promise<Event> {
    let headers = new Headers({
      'Content-Type': 'application/json'});

    return this.http
               .post(this.heroesUrl, JSON.stringify(hero), {headers: headers})
               .toPromise()
               .then(res => res.json())
               .catch(this.handleError);
  }

  // Update existing Event
  private put(hero: Event) {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');

    let url = `${this.heroesUrl}/${hero.id}`;

    return this.http
               .put(url, JSON.stringify(hero), {headers: headers})
               .toPromise()
               .then(() => hero)
               .catch(this.handleError);
  }

  private handleError(error: any) {
    console.error('An error occurred', error);
    return Promise.reject(error.message || error);
  }
}
