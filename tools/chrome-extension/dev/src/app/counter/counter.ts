import { Component, inject, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { counterState, CounterState, increment } from './counter.store';
import { AsyncPipe, JsonPipe } from '@angular/common';

@Component({
  selector: 'app-counter',
  templateUrl: './counter.html',
  styleUrl: './counter.css',
  imports: [AsyncPipe, JsonPipe],
})
export class Counter implements OnInit {
  private readonly store: Store<CounterState> = inject(Store);

  readonly counterState = this.store.select(counterState);

  ngOnInit(): void {
    this.nextFrame();
  }

  private nextFrame() {
    // On every frame, increment the counter
    requestAnimationFrame(() => {
      this.store.dispatch(increment());
      this.nextFrame();
    });
  }
}
