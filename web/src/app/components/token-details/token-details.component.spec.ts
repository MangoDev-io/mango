import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TokenDetailsComponent } from './token-details.component';

describe('TokenDetailsComponent', () => {
  let component: TokenDetailsComponent;
  let fixture: ComponentFixture<TokenDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TokenDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TokenDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
