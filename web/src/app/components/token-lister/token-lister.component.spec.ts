import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TokenListerComponent } from './token-lister.component';

describe('TokenListerComponent', () => {
  let component: TokenListerComponent;
  let fixture: ComponentFixture<TokenListerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TokenListerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TokenListerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
