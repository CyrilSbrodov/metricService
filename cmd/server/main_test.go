package main

//func TestAbs(t *testing.T) {
//
//	tests := []struct {
//		name  string
//		value float64
//		want  float64
//	}{
//		{name: "simple negative value", value: -10, want: 10},
//		{name: "simple positive value", value: 10, want: 10},
//		{name: "zero", value: -0, want: 0},
//		{name: "small value", value: -0.000000001, want: 0.000000001},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			v := Abs(tt.value)
//			// для сравнения двух чисел подойдет функция Equal
//			assert.Equal(t, tt.want, v)
//		})
//	}
//}
//
//func TestUser_FullName(t *testing.T) {
//	type fields struct {
//		FirstName string
//		LastName  string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		{
//			name: "simple test",
//			fields: fields{
//				FirstName: "Misha",
//				LastName:  "Popov",
//			},
//			want: "Misha Popov",
//		},
//		{
//			name: "long name",
//			fields: fields{
//				FirstName: "Pablo Diego KHoze Frantsisko de Paula KHuan Nepomukeno Krispin Krispiano de la Santisima Trinidad Ruiz",
//				LastName:  "Picasso",
//			},
//			want: "Pablo Diego KHoze Frantsisko de Paula KHuan Nepomukeno Krispin Krispiano de la Santisima Trinidad Ruiz Picasso",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := User{
//				FirstName: tt.fields.FirstName,
//				LastName:  tt.fields.LastName,
//			}
//			v := u.FullName()
//			// как и в предыдущем тесте сроки сравниваются с помощью функции Equal
//			assert.Equal(t, tt.want, v)
//		})
//	}
//}
//
//func TestFamily_AddNew(t *testing.T) {
//	type newPerson struct {
//		r Relationship
//		p Person
//	}
//	tests := []struct {
//		name           string
//		existedMembers map[Relationship]Person
//		newPerson      newPerson
//		wantErr        bool
//	}{
//		{
//			name: "add father",
//			existedMembers: map[Relationship]Person{
//				Mother: {
//					FirstName: "Maria",
//					LastName:  "Popova",
//					Age:       36,
//				},
//			},
//			newPerson: newPerson{
//				r: Father,
//				p: Person{
//					FirstName: "Misha",
//					LastName:  "Popov",
//					Age:       42,
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "catch error",
//			existedMembers: map[Relationship]Person{
//				Father: {
//					FirstName: "Misha",
//					LastName:  "Popov",
//					Age:       42,
//				},
//			},
//			newPerson: newPerson{
//				r: Father,
//				p: Person{
//					FirstName: "Ken",
//					LastName:  "Gymsohn",
//					Age:       32,
//				},
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			f := &Family{
//				Members: tt.existedMembers,
//			}
//			err := f.AddNew(tt.newPerson.r, tt.newPerson.p)
//			if !tt.wantErr {
//				// обязательно проверяем на ошибки
//				require.NoError(t, err)
//				// дополнительно проверяем, что новый человек был добавлен
//				assert.Contains(t, f.Members, tt.newPerson.r)
//				return
//			}
//
//			assert.Error(t, err)
//		})
//	}
//}
