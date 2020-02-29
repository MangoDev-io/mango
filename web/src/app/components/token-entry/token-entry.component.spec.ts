import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TokenEntryComponent } from './token-entry.component';

describe('TokenEntryComponent', () => {
  let component: TokenEntryComponent;
  let fixture: ComponentFixture<TokenEntryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TokenEntryComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TokenEntryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
