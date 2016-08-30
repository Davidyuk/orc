import { Component, OnInit } from '@angular/core';
import { Router }            from '@angular/router';

import { Event }                from './hero';
import { EventService }         from './hero.service';

import { Observable }        from 'rxjs/Observable';
import { Subject }           from 'rxjs/Subject';

@Component({
  selector: 'my-heroes',
  templateUrl: 'app/heroes.component.html',
  styleUrls:  ['app/heroes.component.css']
})
export class EventsComponent implements OnInit {
  heroes: Observable<Event[]>;
  private searchTerms = new Subject<string>();
  selectedHero: Event;
  addingHero = false;
  error: any;

  constructor(
    private router: Router,
    private heroService: EventService) { }

  getHeroes() {
    this.heroService
      .getHeroes()
      .then(heroes => this.heroes = Observable.of<Event[]>(heroes))
      .catch(error => this.error = error);
  }

  search(term: string) { this.searchTerms.next(term); }

  addHero() {
    this.addingHero = true;
    this.selectedHero = null;
  }

  close(savedHero: Event) {
    this.addingHero = false;
    alert('if (savedHero) { this.getHeroes(); }');
  }

  deleteHero(hero: Event, event: any) {
    event.stopPropagation();
    this.heroService
        .delete(hero)
        .catch(error => this.error = error);
  }

  ngOnInit() {
    this.heroes = this.searchTerms
      .debounceTime(300).distinctUntilChanged()
      .switchMap(term => this.heroService.search(term))
      .catch(error => {
        this.error = error;
        return Observable.of<Event[]>([]);
      });
    // this.search('');
    setTimeout(() => this.search(''));
  }

  onSelect(hero: Event) {
    this.selectedHero = hero;
    this.addingHero = false;
  }

  gotoDetail(hero: Event) {
    this.router.navigate(['/detail', hero.id]);
  }
}
