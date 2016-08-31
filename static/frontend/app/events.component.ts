import { Component, OnInit } from '@angular/core';
import { Router }            from '@angular/router';

import { Event }                from './event';
import { EventService }         from './event.service';

import { Observable }        from 'rxjs/Observable';
import { Subject }           from 'rxjs/Subject';

@Component({
  selector: 'my-events',
  templateUrl: 'app/events.component.html',
  styleUrls:  ['app/events.component.css']
})
export class EventsComponent implements OnInit {
  events: Observable<Event[]>;
  private searchTerms = new Subject<string>();
  error: any;

  constructor(
    private router: Router,
    private eventService: EventService) { }

  getEvent() {
    this.eventService
      .getEvents()
      .then(events => this.events = Observable.of<Event[]>(events))
      .catch(error => this.error = error);
  }

  search(term: string) { this.searchTerms.next(term); }

  deleteEvent(event: Event, e: any) {
    e.stopPropagation();
    this.eventService
        .delete(event)
        .catch(error => this.error = error);
  }

  ngOnInit() {
    this.events = this.searchTerms
      .debounceTime(300).distinctUntilChanged()
      .switchMap(term => this.eventService.search(term))
      .catch(error => {
        this.error = error;
        return Observable.of<Event[]>([]);
      });
    // this.search('');
    setTimeout(() => this.search(''));
  }

  gotoDetail(event: Event) {
    this.router.navigate(['/detail', event.id]);
  }
}
