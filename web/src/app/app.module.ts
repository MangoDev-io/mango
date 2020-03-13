import { BrowserModule } from '@angular/platform-browser'
import { NgModule } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { HttpClientModule } from '@angular/common/http'

import { AppRoutingModule } from './app-routing.module'
import { AppComponent } from './app.component'
import { LoginComponent } from './pages/login/login.component'
import { ManageComponent } from './pages/manage/manage.component'
import { TokenListerComponent } from './components/token-lister/token-lister.component'
import { TokenEntryComponent } from './components/token-entry/token-entry.component'
import { TokenDetailsComponent } from './components/token-details/token-details.component'
import { TokenCreateComponent } from './components/token-create/token-create.component';
import { NotificationComponent } from './components/notification/notification.component'

@NgModule({
    declarations: [
        AppComponent,
        LoginComponent,
        ManageComponent,
        TokenListerComponent,
        TokenEntryComponent,
        TokenDetailsComponent,
        TokenCreateComponent,
        NotificationComponent,
    ],
    imports: [BrowserModule, AppRoutingModule, FormsModule, HttpClientModule],
    providers: [],
    bootstrap: [AppComponent],
})
export class AppModule {}
