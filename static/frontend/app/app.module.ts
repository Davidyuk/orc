import { NgModule }       from '@angular/core';
import { BrowserModule }  from '@angular/platform-browser';
import { FormsModule }    from '@angular/forms';

import { HttpModule, RequestOptions } from '@angular/http';

import { AppComponent }   from './app.component';
import { routing }        from './app.routing';

import { EventsComponent }      from './heroes.component';
import { EventDetailComponent }  from './hero-detail.component';

import { EventService }  from './hero.service';
import { AppRequestOptions } from "./request-options.service";

import {MdButtonModule} from "@angular2-material/button";
import {MdCardModule} from "@angular2-material/card";
import {MdToolbarModule} from "@angular2-material/toolbar";
import {MdInputModule} from "@angular2-material/input";
import {MdListModule} from "@angular2-material/list";

@NgModule({
  imports: [
    BrowserModule,
    FormsModule,
    routing,
    HttpModule,
    MdCardModule, MdButtonModule, MdToolbarModule, MdInputModule, MdListModule,
  ],
  declarations: [
    AppComponent,
    EventsComponent,
    EventDetailComponent,
  ],
  providers: [
    EventService,
    { provide: RequestOptions, useClass: AppRequestOptions },
    { provide: 'apiBaseUrl', useValue: 'http://localhost:5000/' },
  ],
  bootstrap: [ AppComponent ]
})
export class AppModule {
}


/*
Copyright 2016 Google Inc. All Rights Reserved.
Use of this source code is governed by an MIT-style license that
can be found in the LICENSE file at http://angular.io/license
*/
