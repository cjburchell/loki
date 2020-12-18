import {Component, Input, OnInit} from '@angular/core';
import {IEndpoint} from '../../services/contracts/Endpoint';

@Component({
  selector: 'app-endpoint-item',
  templateUrl: './endpoint-item.component.html',
  styleUrls: ['./endpoint-item.component.scss']
})
export class EndpointItemComponent implements OnInit {
  @Input() endpoint: IEndpoint | undefined;

  ngOnInit(): void {
  }
}
