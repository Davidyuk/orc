import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

import { ActivatedRoute, Params } from '@angular/router';

import { Event }        from './event';
import { EventService } from './event.service';

@Component({
  selector: 'my-event-detail',
  templateUrl: 'app/event-detail.component.html',
  styleUrls: ['app/event-detail.component.css']
})
export class EventDetailComponent implements OnInit {
  @Input() event: Event;
  @Output() close = new EventEmitter();
  error: any;
  navigated = false; // true if navigated here

  constructor(
    private eventService: EventService,
    private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.params.forEach((params: Params) => {
      if (params['id'] !== undefined) {
        let id = +params['id'];
        this.navigated = true;
        this.eventService.getEvent(id)
            .then(event => this.event = event);
      } else {
        this.navigated = false;
        this.event = new Event();
      }
    });
  }

  save() {
    this.eventService
        .save(this.event)
        .then(event => {
          this.event = event; // saved event, w/ id if new
          this.goBack(event);
        })
        .catch(error => this.error = error); // TODO: Display error message
  }
  goBack(savedEvent: Event = null) {
    this.close.emit(savedEvent);
    if (this.navigated) { window.history.back(); }
  }
}
