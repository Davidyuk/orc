import { Routes, RouterModule } from '@angular/router';

import { EventsComponent }     from './heroes.component';
import { EventDetailComponent } from './hero-detail.component';

const appRoutes: Routes = [
  {
    path: '',
    redirectTo: '/events',
    pathMatch: 'full'
  },
  {
    path: 'detail/:id',
    component: EventDetailComponent
  },
  {
    path: 'events',
    component: EventsComponent
  }
];

export const routing = RouterModule.forRoot(appRoutes);
