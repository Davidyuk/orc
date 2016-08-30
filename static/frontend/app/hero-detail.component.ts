import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

import { ActivatedRoute, Params } from '@angular/router';

import { Event }        from './hero';
import { EventService } from './hero.service';

@Component({
  selector: 'my-hero-detail',
  templateUrl: 'app/hero-detail.component.html',
  styleUrls: ['app/hero-detail.component.css']
})
export class EventDetailComponent implements OnInit {
  @Input() hero: Event;
  @Output() close = new EventEmitter();
  error: any;
  navigated = false; // true if navigated here

  constructor(
    private heroService: EventService,
    private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.params.forEach((params: Params) => {
      if (params['id'] !== undefined) {
        let id = +params['id'];
        this.navigated = true;
        this.heroService.getHero(id)
            .then(hero => this.hero = hero);
      } else {
        this.navigated = false;
        this.hero = new Event();
      }
    });
  }

  save() {
    this.heroService
        .save(this.hero)
        .then(hero => {
          this.hero = hero; // saved hero, w/ id if new
          this.goBack(hero);
        })
        .catch(error => this.error = error); // TODO: Display error message
  }
  goBack(savedHero: Event = null) {
    this.close.emit(savedHero);
    if (this.navigated) { window.history.back(); }
  }
}


/*
Copyright 2016 Google Inc. All Rights Reserved.
Use of this source code is governed by an MIT-style license that
can be found in the LICENSE file at http://angular.io/license
*/
