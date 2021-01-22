import {Component, Input, EventEmitter, Output} from '@angular/core';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent {
  public searchValue: string | undefined;

  @Input()
  get searchText(): string {
    return this.searchValue as string;
  }

  set searchText(val) {
    this.searchValue = val;
    this.searchTextChange.emit(val);
  }

  @Output()
  searchTextChange = new EventEmitter<string>();

}
